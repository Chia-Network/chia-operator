{{/*
Expand the name of the chart.
*/}}
{{- define "chia.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "chia.fullname" -}}
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
{{- define "chia.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "chia.labels" -}}
helm.sh/chart: {{ include "chia.chart" . }}
{{ include "chia.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "chia.selectorLabels" -}}
app.kubernetes.io/name: {{ include "chia.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "chia.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "chia.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the CA secret to use
*/}}
{{- define "chia.caSecretName" -}}
{{- .Values.ca.secretName | default (printf "%s-ca" (include "chia.fullname" .)) }}
{{- end }}

{{/*
Common Chia configuration block
*/}}
{{- define "chia.commonConfig" -}}
{{- if .Values.testnet }}
testnet: {{ .Values.testnet }}
{{- end }}
timezone: {{ .Values.timezone }}
logLevel: {{ .Values.logLevel }}
selfHostname: {{ .Values.selfHostname }}
{{- if .Values.sourceRef }}
sourceRef: {{ .Values.sourceRef }}
{{- end }}
{{- end -}}

{{/*
Chia Exporter configuration block
*/}}
{{- define "chia.exporterConfig" -}}
{{- if .Values.chiaExporter.enabled }}
chiaExporter:
  enabled: {{ .Values.chiaExporter.enabled }}
  {{- if .Values.chiaExporter.image }}
  image: {{ .Values.chiaExporter.image }}
  {{- end }}
  service:
    {{- if .Values.chiaExporter.service.annotations }}
    annotations: {{ toYaml .Values.chiaExporter.service.annotations | nindent 6 }}
    {{- end }}
    {{- if .Values.chiaExporter.service.labels }}
    labels: {{ toYaml .Values.chiaExporter.service.labels | nindent 6 }}
    {{- end }}
{{- end }}
{{- end -}}
