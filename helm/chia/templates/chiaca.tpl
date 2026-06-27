{{- if or (empty .Values.ca.chia_ca_crt) (empty .Values.ca.chia_ca_key) (empty .Values.ca.private_ca_crt) (empty .Values.ca.private_ca_key) -}}
apiVersion: k8s.chia.net/v1
kind: ChiaCA
metadata:
  name: {{ include "chia.fullname" . }}
spec:
  secret: {{ include "chia.caSecretName" . }}
{{- else }}
apiVersion: v1
kind: Secret
metadata:
  name: {{ include "chia.caSecretName" . }}
stringData:
  chia_ca.crt: |
    {{ .Values.ca.chia_ca_crt | nindent 4 }}
  chia_ca.key: |
    {{ .Values.ca.chia_ca_key | nindent 4 }}
  private_ca.crt: |
    {{ .Values.ca.private_ca_crt | nindent 4 }}
  private_ca.key: |
    {{ .Values.ca.private_ca_key | nindent 4 }}
type: Opaque
{{- end }}
