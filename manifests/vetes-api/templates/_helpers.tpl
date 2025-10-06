{{/*
Expand the name of the chart.
*/}}
{{- define "vetes-api.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "vetes-api.fullname" -}}
{{- if contains .Chart.Name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "vetes-api.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "vetes-api.labels" -}}
helm.sh/chart: {{ include "vetes-api.chart" . }}
{{ include "vetes-api.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Values.labels }}
{{ toYaml .Values.labels }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "vetes-api.selectorLabels" -}}
app.kubernetes.io/name: {{ include "vetes-api.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Return the proper image name
*/}}
{{- define "vetes-api.image" -}}
{{- $repositoryName := .Values.platformConfig.imageRepositoryRelease -}}
{{- $imageName := .Values.image.name -}}
{{- $tag := (default .Chart.AppVersion .Values.image.tag) | toString -}}
{{- if .Values.platformConfig.imageRegistry -}}
{{- $registryName := .Values.platformConfig.imageRegistry -}}
{{- printf "%s/%s/%s:%s" $registryName $repositoryName $imageName $tag -}}
{{- else -}}
{{- printf "%s/%s:%s" $repositoryName $imageName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Create the normalize configmap name
*/}}
{{- define "vetes-api.normalizeName" -}}
{{ include "vetes-api.fullname" . | trunc 53 | trimSuffix "-" }}-normalize
{{- end -}}
