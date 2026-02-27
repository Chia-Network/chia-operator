{{- if .Values.node.enabled -}}
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: {{ include "chia.fullname" . }}
spec:
  replicas: {{ .Values.node.replicas }}
  chia:
    caSecretName: {{ include "chia.caSecretName" . }}
    {{- if gt (len .Values.node.fullNodePeers) 0 }}
    fullNodePeers: {{ toYaml .Values.node.fullNodePeers | nindent 6 }}
    {{- end }}
    {{ include "chia.commonConfig" . | nindent 4 }}
  {{ include "chia.exporterConfig" . | nindent 2 }}
  {{ include "chia.healthcheckConfig" . | nindent 2 }}
  {{- if .Values.node.storage.chiaRoot.persistentVolumeClaim.enabled }}
  storage:
    chiaRoot:
      persistentVolumeClaim:
        storageClass: "{{ .Values.node.storage.chiaRoot.persistentVolumeClaim.storageClass }}"
        resourceRequest: "{{ .Values.node.storage.chiaRoot.persistentVolumeClaim.resourceRequest }}"
  {{- end }}
{{- end }}
