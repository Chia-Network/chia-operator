{{- if .Values.node.enabled -}}
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: {{ include "chia.fullname" . }}
spec:
  chia:
    caSecretName: {{ include "chia.caSecretName" . }}
{{- end }}
