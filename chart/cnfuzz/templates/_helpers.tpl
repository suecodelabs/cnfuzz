{{/*
Expand the name of the chart.
*/}}
{{- define "cnfuzz.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cnfuzz.fullname" -}}
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
{{- define "cnfuzz.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cnfuzz.labels" -}}
helm.sh/chart: {{ include "cnfuzz.chart" . }}
{{ include "cnfuzz.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cnfuzz.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cnfuzz.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "cnfuzz.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "cnfuzz.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use for restlerwrapper jobs
*/}}
{{- define "restlerwrapper.serviceAccountName" -}}
{{- if .Values.restlerwrapper.serviceAccount.create }}
{{- default (printf "%s-job" (include "cnfuzz.fullname" .)) .Values.restlerwrapper.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.restlerwrapper.serviceAccount.name }}
{{- end }}
{{- end }}
