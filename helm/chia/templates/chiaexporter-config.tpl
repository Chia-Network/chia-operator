{{- if .Values.chiaExporter.additionalConfig }}
apiVersion: v1
kind: Secret
metadata:
  name: {{ include "chia.exporterSecretName" . }}
stringData:
  {{ toYaml .Values.chiaExporter.additionalConfig | nindent 2 }}
{{- end }}
