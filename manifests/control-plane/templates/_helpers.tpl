{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "controlPlane.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "controlPlane.fullname" -}}
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
{{- define "controlPlane.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "controlPlane.labels" -}}
helm.sh/chart: {{ include "controlPlane.chart" . }}
{{ include "controlPlane.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "controlPlane.selectorLabels" -}}
app.kubernetes.io/name: {{ include "controlPlane.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Name of ConfigMap containing pipecd configuration
*/}}
{{- define "controlPlane.configMapName" -}}
{{- if .Values.config.create }}
{{- include "controlPlane.fullname" . }}
{{- else }}
{{- .Values.config.name }}
{{- end }}
{{- end }}

{{/*
Name of Secret containing sensitive data
*/}}
{{- define "controlPlane.secretName" -}}
{{- if .Values.secret.create }}
{{- include "controlPlane.fullname" . }}
{{- else }}
{{- .Values.secret.name }}
{{- end }}
{{- end }}

{{/*
Name of ServiceAccount
*/}}
{{- define "controlPlane.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{ include "controlPlane.fullname" . }}
{{- else }}
{{- .Values.serviceAccount.name }}
{{- end }}
{{- end }}
