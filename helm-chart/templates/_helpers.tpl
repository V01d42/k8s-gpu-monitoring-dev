{{/*
Expand the name of the chart.
*/}}
{{- define "k8s-gpu-monitoring-dev.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "k8s-gpu-monitoring-dev.fullname" -}}
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
{{- define "k8s-gpu-monitoring-dev.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "k8s-gpu-monitoring-dev.labels" -}}
helm.sh/chart: {{ include "k8s-gpu-monitoring-dev.chart" . }}
{{ include "k8s-gpu-monitoring-dev.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- with .Values.commonLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "k8s-gpu-monitoring-dev.selectorLabels" -}}
app.kubernetes.io/name: {{ include "k8s-gpu-monitoring-dev.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Backend specific labels
*/}}
{{- define "k8s-gpu-monitoring-dev.backend.labels" -}}
{{ include "k8s-gpu-monitoring-dev.labels" . }}
app.kubernetes.io/component: backend
{{- with .Values.backend.podLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Backend selector labels
*/}}
{{- define "k8s-gpu-monitoring-dev.backend.selectorLabels" -}}
{{ include "k8s-gpu-monitoring-dev.selectorLabels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Backend fullname
*/}}
{{- define "k8s-gpu-monitoring-dev.backend.fullname" -}}
{{ include "k8s-gpu-monitoring-dev.fullname" . }}-backend
{{- end }}

{{/*
Frontend specific labels
*/}}
{{- define "k8s-gpu-monitoring-dev.frontend.labels" -}}
{{ include "k8s-gpu-monitoring-dev.labels" . }}
app.kubernetes.io/component: frontend
{{- with .Values.frontend.podLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Frontend selector labels
*/}}
{{- define "k8s-gpu-monitoring-dev.frontend.selectorLabels" -}}
{{ include "k8s-gpu-monitoring-dev.selectorLabels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Frontend fullname
*/}}
{{- define "k8s-gpu-monitoring-dev.frontend.fullname" -}}
{{ include "k8s-gpu-monitoring-dev.fullname" . }}-frontend
{{- end }}

{{/*
Backend image name
*/}}
{{- define "k8s-gpu-monitoring-dev.backend.image" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s:%s" .Values.global.imageRegistry .Values.backend.image.repository .Values.backend.image.tag }}
{{- else }}
{{- printf "%s:%s" .Values.backend.image.repository .Values.backend.image.tag }}
{{- end }}
{{- end }}

{{/*
Frontend image name
*/}}
{{- define "k8s-gpu-monitoring-dev.frontend.image" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s:%s" .Values.global.imageRegistry .Values.frontend.image.repository .Values.frontend.image.tag }}
{{- else }}
{{- printf "%s:%s" .Values.frontend.image.repository .Values.frontend.image.tag }}
{{- end }}
{{- end }}

{{/*
Common annotations
*/}}
{{- define "k8s-gpu-monitoring-dev.annotations" -}}
{{- with .Values.commonAnnotations }}
{{ toYaml . }}
{{- end }}
{{- end }} 