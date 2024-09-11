{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "piped.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "piped.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "piped.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "piped.labels" -}}
helm.sh/chart: {{ include "piped.chart" . }}
{{ include "piped.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "piped.selectorLabels" -}}
app.kubernetes.io/name: {{ include "piped.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Name of ConfigMap containing piped configuration
*/}}
{{- define "piped.configMapName" -}}
{{- if .Values.config.create }}
{{- include "piped.fullname" . }}
{{- else }}
{{- .Values.config.name }}
{{- end }}
{{- end }}

{{/*
Name of Secret containing sensitive data
*/}}
{{- define "piped.secretName" -}}
{{- if .Values.secret.create }}
{{- include "piped.fullname" . }}
{{- else }}
{{- .Values.secret.name }}
{{- end }}
{{- end }}

{{/*
Name of ServiceAccount
*/}}
{{- define "piped.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{ include "piped.fullname" . }}
{{- else }}
{{- .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
A set of permissions Role will contain
*/}}
{{- define "piped.roleRules" -}}
{{- if .Values.rbac.rules -}}
{{- with .Values.rbac.rules }}
{{- toYaml . | nindent 2 }}
{{- end }}
{{- else }}
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
{{- end }}
{{- end }}

{{/*
A set of permissions ClusterRole will contain
*/}}
{{- define "piped.clusterRoleRules" -}}
{{- if .Values.rbac.rules -}}
{{- with .Values.rbac.rules }}
{{- toYaml . | nindent 2 }}
{{- end }}
{{- else }}
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs:
  - '*'
  verbs:
  - '*'
{{- end }}
{{- end }}

{{/*
A set of args for Launcher.
*/}}
{{- define "piped.launcherArgs" -}}
- launcher
{{- if .Values.launcher.configFromGitRepo.enabled }}
{{- with .Values.launcher.configFromGitRepo }}
- --config-from-git-repo=true
- --git-repo-url={{ required "repoUrl is required" .repoUrl }}
- --git-branch={{ required "branch is required" .branch }}
- --git-piped-config-file={{ required "configFile is required" .configFile }}
- --git-ssh-key-file={{ required "sshKeyFile is required" .sshKeyFile }}
{{- end }}
{{- else }}
- --config-file=/etc/piped-config/{{ .Values.config.fileName }}
{{- end }}
- --metrics={{ .Values.args.metrics }}
- --enable-default-kubernetes-cloud-provider={{ .Values.args.enableDefaultKubernetesCloudProvider }}
- --log-encoding={{ .Values.args.logEncoding }}
- --log-level={{ .Values.args.logLevel }}
- --add-login-user-to-passwd={{ .Values.args.addLoginUserToPasswd }}
{{- if .Values.quickstart.enabled }}
- --insecure=true
{{- else }}
- --insecure={{ .Values.args.insecure }}
{{- end }}
{{- end }}

{{/*
A set of args for Piped.
*/}}
{{- define "piped.pipedArgs" -}}
- piped
- --config-file=/etc/piped-config/{{ .Values.config.fileName }}
- --metrics={{ .Values.args.metrics }}
- --enable-default-kubernetes-cloud-provider={{ .Values.args.enableDefaultKubernetesCloudProvider }}
- --log-encoding={{ .Values.args.logEncoding }}
- --log-level={{ .Values.args.logLevel }}
- --add-login-user-to-passwd={{ .Values.args.addLoginUserToPasswd }}
- --app-manifest-cache-count={{ .Values.args.appManifestCacheCount }}
{{- if .Values.quickstart.enabled }}
- --insecure=true
{{- else }}
- --insecure={{ .Values.args.insecure }}
{{- end }}
{{- end }}
