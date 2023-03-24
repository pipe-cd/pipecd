// Copyright 2023 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/yamlprocessor"
)

type deployExecutor struct {
	executor.Input

	commit string
	appCfg *config.KubernetesApplicationSpec

	loader        provider.Loader
	applierGetter applierGetter
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.RollbackKind, f executor.Factory) error
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &deployExecutor{
			Input: in,
		}
	}

	r.Register(model.StageK8sSync, f)
	r.Register(model.StageK8sPrimaryRollout, f)
	r.Register(model.StageK8sCanaryRollout, f)
	r.Register(model.StageK8sCanaryClean, f)
	r.Register(model.StageK8sBaselineRollout, f)
	r.Register(model.StageK8sBaselineClean, f)
	r.Register(model.StageK8sTrafficRouting, f)

	r.RegisterRollback(model.RollbackKind_Rollback_KUBERNETES, func(in executor.Input) executor.Executor {
		return &rollbackExecutor{
			Input: in,
		}
	})
}

func (e *deployExecutor) Execute(sig executor.StopSignal) model.StageStatus {
	ctx := sig.Context()
	e.commit = e.Deployment.Trigger.Commit.Hash

	ds, err := e.TargetDSP.Get(ctx, e.LogPersister)
	if err != nil {
		e.LogPersister.Errorf("Failed to prepare target deploy source data (%v)", err)
		return model.StageStatus_STAGE_FAILURE
	}

	e.appCfg = ds.ApplicationConfig.KubernetesApplicationSpec
	if e.appCfg == nil {
		e.LogPersister.Error("Malformed application configuration: missing KubernetesApplicationSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	if e.appCfg.Input.HelmChart != nil {
		chartRepoName := e.appCfg.Input.HelmChart.Repository
		if chartRepoName != "" {
			e.appCfg.Input.HelmChart.Insecure = e.PipedConfig.IsInsecureChartRepository(chartRepoName)
		}
	}

	if e.appCfg.Input.KubectlVersion != "" {
		e.LogPersister.Infof("kubectl version %s will be used.", e.appCfg.Input.KubectlVersion)
	}

	e.applierGetter, err = newApplierGroup(e.Deployment.PlatformProvider, *e.appCfg, e.PipedConfig, e.Logger)
	if err != nil {
		e.LogPersister.Error(err.Error())
		return model.StageStatus_STAGE_FAILURE
	}

	e.loader = provider.NewLoader(
		e.Deployment.ApplicationName,
		ds.AppDir,
		ds.RepoDir,
		e.Deployment.GitPath.ConfigFilename,
		e.appCfg.Input,
		e.GitClient,
		e.Logger,
	)

	e.Logger.Info("start executing kubernetes stage",
		zap.String("stage-name", e.Stage.Name),
		zap.String("app-dir", ds.AppDir),
	)

	var (
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageK8sSync:
		status = e.ensureSync(ctx)

	case model.StageK8sPrimaryRollout:
		status = e.ensurePrimaryRollout(ctx)

	case model.StageK8sCanaryRollout:
		status = e.ensureCanaryRollout(ctx)

	case model.StageK8sCanaryClean:
		status = e.ensureCanaryClean(ctx)

	case model.StageK8sBaselineRollout:
		status = e.ensureBaselineRollout(ctx)

	case model.StageK8sBaselineClean:
		status = e.ensureBaselineClean(ctx)

	case model.StageK8sTrafficRouting:
		status = e.ensureTrafficRouting(ctx)

	default:
		e.LogPersister.Errorf("Unsupported stage %s for kubernetes application", e.Stage.Name)
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *deployExecutor) loadRunningManifests(ctx context.Context) (manifests []provider.Manifest, err error) {
	commit := e.Deployment.RunningCommitHash
	if commit == "" {
		return nil, fmt.Errorf("unable to determine running commit")
	}

	loader := &manifestsLoadFunc{
		loadFunc: func(ctx context.Context) ([]provider.Manifest, error) {
			ds, err := e.RunningDSP.Get(ctx, e.LogPersister)
			if err != nil {
				e.LogPersister.Errorf("Failed to prepare running deploy source (%v)", err)
				return nil, err
			}

			loader := provider.NewLoader(
				e.Deployment.ApplicationName,
				ds.AppDir,
				ds.RepoDir,
				e.Deployment.GitPath.ConfigFilename,
				e.appCfg.Input,
				e.GitClient,
				e.Logger,
			)
			return loader.LoadManifests(ctx)
		},
	}

	return loadManifests(ctx, e.Deployment.ApplicationId, commit, e.AppManifestsCache, loader, e.Logger)
}

type manifestsLoadFunc struct {
	loadFunc func(context.Context) ([]provider.Manifest, error)
}

func (l *manifestsLoadFunc) LoadManifests(ctx context.Context) ([]provider.Manifest, error) {
	return l.loadFunc(ctx)
}

func loadManifests(ctx context.Context, appID, commit string, manifestsCache cache.Cache, loader provider.Loader, logger *zap.Logger) (manifests []provider.Manifest, err error) {
	cache := provider.AppManifestsCache{
		AppID:  appID,
		Cache:  manifestsCache,
		Logger: logger,
	}
	manifests, ok := cache.Get(commit)
	if ok {
		return manifests, nil
	}

	// When the manifests were not in the cache we have to load them.
	if manifests, err = loader.LoadManifests(ctx); err != nil {
		return nil, err
	}
	cache.Put(commit, manifests)

	return manifests, nil
}

func addBuiltinAnnotations(manifests []provider.Manifest, variantLabel, variant, hash, pipedID, appID string) {
	for i := range manifests {
		manifests[i].AddAnnotations(map[string]string{
			provider.LabelManagedBy:          provider.ManagedByPiped,
			provider.LabelPiped:              pipedID,
			provider.LabelApplication:        appID,
			variantLabel:                     variant,
			provider.LabelOriginalAPIVersion: manifests[i].Key.APIVersion,
			provider.LabelResourceKey:        manifests[i].Key.String(),
			provider.LabelCommitHash:         hash,
		})
	}
}

func applyManifests(ctx context.Context, ag applierGetter, manifests []provider.Manifest, namespace string, lp executor.LogPersister) error {
	if namespace == "" {
		lp.Infof("Start applying %d manifests", len(manifests))
	} else {
		lp.Infof("Start applying %d manifests to %q namespace", len(manifests), namespace)
	}

	for _, m := range manifests {
		applier, err := ag.Get(m.Key)
		if err != nil {
			lp.Error(err.Error())
			return err
		}

		annotation := m.GetAnnotations()[provider.LabelSyncReplace]
		if annotation != provider.UseReplaceEnabled {
			if err := applier.ApplyManifest(ctx, m); err != nil {
				lp.Errorf("Failed to apply manifest: %s (%w)", m.Key.ReadableString(), err)
				return err
			}
			lp.Successf("- applied manifest: %s", m.Key.ReadableString())
			continue
		}
		// Always try to replace first and create if it fails due to resource not found error.
		// This is because we cannot know whether resource already exists before executing command.
		err = applier.ReplaceManifest(ctx, m)
		if errors.Is(err, provider.ErrNotFound) {
			lp.Infof("Specified resource does not exist, so create the resource: %s (%w)", m.Key.ReadableString(), err)
			err = applier.CreateManifest(ctx, m)
		}
		if err != nil {
			lp.Errorf("Failed to replace or create manifest: %s (%w)", m.Key.ReadableString(), err)
			return err
		}
		lp.Successf("- replaced or created manifest: %s", m.Key.ReadableString())

	}
	lp.Successf("Successfully applied %d manifests", len(manifests))
	return nil
}

func deleteResources(ctx context.Context, ag applierGetter, resources []provider.ResourceKey, lp executor.LogPersister) error {
	resourcesLen := len(resources)
	if resourcesLen == 0 {
		lp.Info("No resources to delete")
		return nil
	}

	lp.Infof("Start deleting %d resources", len(resources))
	var deletedCount int

	for _, k := range resources {
		applier, err := ag.Get(k)
		if err != nil {
			lp.Error(err.Error())
			return err
		}

		err = applier.Delete(ctx, k)
		if err == nil {
			lp.Successf("- deleted resource: %s", k.ReadableString())
			deletedCount++
			continue
		}
		if errors.Is(err, provider.ErrNotFound) {
			lp.Infof("- no resource %s to delete", k.ReadableString())
			deletedCount++
			continue
		}
		lp.Errorf("- unable to delete resource: %s (%v)", k.ReadableString(), err)
	}

	if deletedCount < resourcesLen {
		lp.Infof("Deleted %d/%d resources", deletedCount, resourcesLen)
		return fmt.Errorf("unable to delete %d resources", resourcesLen-deletedCount)
	}

	lp.Successf("Successfully deleted %d resources", len(resources))
	return nil
}

func findManifests(kind, name string, manifests []provider.Manifest) []provider.Manifest {
	var out []provider.Manifest
	for _, m := range manifests {
		if m.Key.Kind != kind {
			continue
		}
		if name != "" && m.Key.Name != name {
			continue
		}
		out = append(out, m)
	}
	return out
}

func findConfigMapManifests(manifests []provider.Manifest) []provider.Manifest {
	var out []provider.Manifest
	for _, m := range manifests {
		if !m.Key.IsConfigMap() {
			continue
		}
		out = append(out, m)
	}
	return out
}

func findSecretManifests(manifests []provider.Manifest) []provider.Manifest {
	var out []provider.Manifest
	for _, m := range manifests {
		if !m.Key.IsSecret() {
			continue
		}
		out = append(out, m)
	}
	return out
}

func findWorkloadManifests(manifests []provider.Manifest, refs []config.K8sResourceReference) []provider.Manifest {
	if len(refs) == 0 {
		return findManifests(provider.KindDeployment, "", manifests)
	}

	workloads := make([]provider.Manifest, 0)
	for _, ref := range refs {
		kind := provider.KindDeployment
		if ref.Kind != "" {
			kind = ref.Kind
		}
		ms := findManifests(kind, ref.Name, manifests)
		workloads = append(workloads, ms...)
	}
	return workloads
}

func duplicateManifests(manifests []provider.Manifest, nameSuffix string) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		out = append(out, duplicateManifest(m, nameSuffix))
	}
	return out
}

func duplicateManifest(m provider.Manifest, nameSuffix string) provider.Manifest {
	name := makeSuffixedName(m.Key.Name, nameSuffix)
	return m.Duplicate(name)
}

func generateVariantServiceManifests(services []provider.Manifest, variantLabel, variant, nameSuffix string) ([]provider.Manifest, error) {
	manifests := make([]provider.Manifest, 0, len(services))
	updateService := func(s *corev1.Service) {
		s.Name = makeSuffixedName(s.Name, nameSuffix)
		// Currently, we suppose that all generated services should be ClusterIP.
		s.Spec.Type = corev1.ServiceTypeClusterIP
		// Append the variant label to the selector
		// to ensure that the generated service is using only workloads of this variant.
		if s.Spec.Selector == nil {
			s.Spec.Selector = map[string]string{}
		}
		s.Spec.Selector[variantLabel] = variant
		// Empty all unneeded fields.
		s.Spec.ExternalIPs = nil
		s.Spec.LoadBalancerIP = ""
		s.Spec.LoadBalancerSourceRanges = nil
	}

	for _, m := range services {
		s := &corev1.Service{}
		if err := m.ConvertToStructuredObject(s); err != nil {
			return nil, err
		}
		updateService(s)
		manifest, err := provider.ParseFromStructuredObject(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Service object to Manifest: %w", err)
		}
		manifests = append(manifests, manifest)
	}
	return manifests, nil
}

func generateVariantWorkloadManifests(workloads, configmaps, secrets []provider.Manifest, variantLabel, variant, nameSuffix string, replicasCalculator func(*int32) int32) ([]provider.Manifest, error) {
	manifests := make([]provider.Manifest, 0, len(workloads))

	cmNames := make(map[string]struct{}, len(configmaps))
	for i := range configmaps {
		cmNames[configmaps[i].Key.Name] = struct{}{}
	}

	secretNames := make(map[string]struct{}, len(secrets))
	for i := range secrets {
		secretNames[secrets[i].Key.Name] = struct{}{}
	}

	updateContainers := func(containers []corev1.Container) {
		for _, container := range containers {
			for _, env := range container.Env {
				if v := env.ValueFrom; v != nil {
					if ref := v.ConfigMapKeyRef; ref != nil {
						if _, ok := cmNames[ref.Name]; ok {
							ref.Name = makeSuffixedName(ref.Name, nameSuffix)
						}
					}
					if ref := v.SecretKeyRef; ref != nil {
						if _, ok := secretNames[ref.Name]; ok {
							ref.Name = makeSuffixedName(ref.Name, nameSuffix)
						}
					}
				}
			}
			for _, envFrom := range container.EnvFrom {
				if ref := envFrom.ConfigMapRef; ref != nil {
					if _, ok := cmNames[ref.Name]; ok {
						ref.Name = makeSuffixedName(ref.Name, nameSuffix)
					}
				}
				if ref := envFrom.SecretRef; ref != nil {
					if _, ok := secretNames[ref.Name]; ok {
						ref.Name = makeSuffixedName(ref.Name, nameSuffix)
					}
				}
			}
		}
	}

	updatePod := func(pod *corev1.PodTemplateSpec) {
		// Add variant labels.
		if pod.Labels == nil {
			pod.Labels = map[string]string{}
		}
		pod.Labels[variantLabel] = variant

		// Update volumes to use canary's ConfigMaps and Secrets.
		for i := range pod.Spec.Volumes {
			if cm := pod.Spec.Volumes[i].ConfigMap; cm != nil {
				if _, ok := cmNames[cm.Name]; ok {
					cm.Name = makeSuffixedName(cm.Name, nameSuffix)
				}
			}
			if s := pod.Spec.Volumes[i].Secret; s != nil {
				if _, ok := secretNames[s.SecretName]; ok {
					s.SecretName = makeSuffixedName(s.SecretName, nameSuffix)
				}
			}
		}

		// Update ENV references in containers.
		updateContainers(pod.Spec.InitContainers)
		updateContainers(pod.Spec.Containers)
	}

	updateDeployment := func(d *appsv1.Deployment) {
		d.Name = makeSuffixedName(d.Name, nameSuffix)
		if replicasCalculator != nil {
			replicas := replicasCalculator(d.Spec.Replicas)
			d.Spec.Replicas = &replicas
		}
		d.Spec.Selector = metav1.AddLabelToSelector(d.Spec.Selector, variantLabel, variant)
		updatePod(&d.Spec.Template)
	}

	for _, m := range workloads {
		switch m.Key.Kind {
		case provider.KindDeployment:
			d := &appsv1.Deployment{}
			if err := m.ConvertToStructuredObject(d); err != nil {
				return nil, err
			}
			updateDeployment(d)
			manifest, err := provider.ParseFromStructuredObject(d)
			if err != nil {
				return nil, err
			}
			manifests = append(manifests, manifest)

		default:
			return nil, fmt.Errorf("unsupported workload kind %s", m.Key.Kind)
		}
	}

	return manifests, nil
}

func checkVariantSelectorInWorkload(m provider.Manifest, variantLabel, variant string) error {
	var (
		matchLabelsFields = []string{"spec", "selector", "matchLabels"}
		labelsFields      = []string{"spec", "template", "metadata", "labels"}
	)

	matchLabels, err := m.GetNestedStringMap(matchLabelsFields...)
	if err != nil {
		return err
	}
	value, ok := matchLabels[variantLabel]
	if !ok {
		return fmt.Errorf("missing %s key in spec.selector.matchLabels", variantLabel)
	}
	if value != variant {
		return fmt.Errorf("require %s but got %s for %s key in %s", variant, value, variantLabel, strings.Join(matchLabelsFields, "."))
	}

	labels, err := m.GetNestedStringMap(labelsFields...)
	if err != nil {
		return err
	}
	value, ok = labels[variantLabel]
	if !ok {
		return fmt.Errorf("missing %s key in spec.template.metadata.labels", variantLabel)
	}
	if value != variant {
		return fmt.Errorf("require %s but got %s for %s key in %s", variant, value, variantLabel, strings.Join(labelsFields, "."))
	}

	return nil
}

func ensureVariantSelectorInWorkload(m provider.Manifest, variantLabel, variant string) error {
	variantMap := map[string]string{
		variantLabel: variant,
	}
	if err := m.AddStringMapValues(variantMap, "spec", "selector", "matchLabels"); err != nil {
		return err
	}
	return m.AddStringMapValues(variantMap, "spec", "template", "metadata", "labels")
}

func makeSuffixedName(name, suffix string) string {
	if suffix != "" {
		return name + "-" + suffix
	}
	return name
}

// annotateConfigHash appends a hash annotation into the workload manifests.
// The hash value is calculated by hashing the content of all configmaps/secrets
// that are referenced by the workload.
// This appending ensures that the workload should be restarted when
// one of its configurations changed.
func annotateConfigHash(manifests []provider.Manifest) error {
	if len(manifests) == 0 {
		return nil
	}

	configMaps := make(map[string]provider.Manifest)
	secrets := make(map[string]provider.Manifest)
	for _, m := range manifests {
		if m.Key.IsConfigMap() {
			configMaps[m.Key.Name] = m
			continue
		}
		if m.Key.IsSecret() {
			secrets[m.Key.Name] = m
		}
	}

	// This application is not containing any config manifests
	// so nothing to do.
	if len(configMaps)+len(secrets) == 0 {
		return nil
	}

	for _, m := range manifests {
		if m.Key.IsDeployment() {
			if err := annotateConfigHashToDeployment(m, configMaps, secrets); err != nil {
				return err
			}
		}

		// TODO: Anotate config hash into other workload kinds such as DaemonSet, StatefulSet...
	}

	return nil
}

func annotateConfigHashToDeployment(m provider.Manifest, managedConfigMaps, managedSecrets map[string]provider.Manifest) error {
	d := &appsv1.Deployment{}
	if err := m.ConvertToStructuredObject(d); err != nil {
		return err
	}

	configMaps := provider.FindReferencingConfigMapsInDeployment(d)
	secrets := provider.FindReferencingSecretsInDeployment(d)

	// The deployment is not referencing any config resources.
	if len(configMaps)+len(secrets) == 0 {
		return nil
	}

	cfgs := make([]provider.Manifest, 0, len(configMaps)+len(secrets))
	for _, cm := range configMaps {
		m, ok := managedConfigMaps[cm]
		if !ok {
			// We do not return error here because the deployment may use
			// a config resource that is not managed by PipeCD.
			continue
		}
		cfgs = append(cfgs, m)
	}
	for _, s := range secrets {
		m, ok := managedSecrets[s]
		if !ok {
			// We do not return error here because the deployment may use
			// a config resource that is not managed by PipeCD.
			continue
		}
		cfgs = append(cfgs, m)
	}

	if len(cfgs) == 0 {
		return nil
	}

	hash, err := provider.HashManifests(cfgs)
	if err != nil {
		return err
	}

	m.AddStringMapValues(
		map[string]string{
			provider.AnnotationConfigHash: hash,
		},
		"spec",
		"template",
		"metadata",
		"annotations",
	)
	return nil
}

type patcher func(m provider.Manifest, cfg config.K8sResourcePatch) (*provider.Manifest, error)

func patchManifests(manifests []provider.Manifest, patches []config.K8sResourcePatch, patcher patcher) ([]provider.Manifest, error) {
	if len(patches) == 0 {
		return manifests, nil
	}

	out := make([]provider.Manifest, len(manifests))
	copy(out, manifests)

	for _, p := range patches {
		target := -1
		for i, m := range out {
			if m.Key.Kind != p.Target.Kind {
				continue
			}
			if m.Key.Name != p.Target.Name {
				continue
			}
			target = i
			break
		}
		if target < 0 {
			return nil, fmt.Errorf("no manifest matches the given patch: kind=%s, name=%s", p.Target.Kind, p.Target.Name)
		}
		patched, err := patcher(out[target], p)
		if err != nil {
			return nil, fmt.Errorf("failed to patch manifest: %s, error: %w", out[target].Key, err)
		}
		out[target] = *patched
	}

	return out, nil
}

func patchManifest(m provider.Manifest, patch config.K8sResourcePatch) (*provider.Manifest, error) {
	if len(patch.Ops) == 0 {
		return &m, nil
	}

	fullBytes, err := m.YamlBytes()
	if err != nil {
		return nil, err
	}

	process := func(bytes []byte) ([]byte, error) {
		p, err := yamlprocessor.NewProcessor(bytes)
		if err != nil {
			return nil, err
		}

		for _, o := range patch.Ops {
			switch o.Op {
			case config.K8sResourcePatchOpYAMLReplace:
				if err := p.ReplaceString(o.Path, o.Value); err != nil {
					return nil, fmt.Errorf("failed to replace value at path: %s, error: %w", o.Path, err)
				}
			default:
				// TODO: Support more patch operation for K8sCanaryRolloutStageOptions.
				return nil, fmt.Errorf("%s operation is not supported currently", o.Op)
			}
		}

		return p.Bytes(), nil
	}

	buildManifest := func(bytes []byte) (*provider.Manifest, error) {
		manifests, err := provider.ParseManifests(string(bytes))
		if err != nil {
			return nil, err
		}
		if len(manifests) != 1 {
			return nil, fmt.Errorf("unexpected number of manifests, expected 1, got %d", len(manifests))
		}
		return &manifests[0], nil
	}

	// When the target is the whole manifest,
	// just pass full bytes to process and build a new manifest based on the returned data.
	root := patch.Target.DocumentRoot
	if root == "" {
		out, err := process(fullBytes)
		if err != nil {
			return nil, err
		}
		return buildManifest(out)
	}

	// When the target is a manifest field specified by documentRoot,
	// we have to extract that field value as a string.
	p, err := yamlprocessor.NewProcessor(fullBytes)
	if err != nil {
		return nil, err
	}

	v, err := p.GetValue(root)
	if err != nil {
		return nil, err
	}
	sv, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("the value for the specified root %s must be a string", root)
	}

	// And process that field data.
	out, err := process([]byte(sv))
	if err != nil {
		return nil, err
	}

	// Then rewrite the new data into the specified root.
	if err := p.ReplaceString(root, string(out)); err != nil {
		return nil, err
	}

	return buildManifest(p.Bytes())
}
