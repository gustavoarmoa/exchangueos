{{- define "exchangeos.fullname" -}}
{{- printf "%s-%s" .Release.Name .binaryName | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "exchangeos.labels" -}}
app.kubernetes.io/name: exchangeos
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: {{ .binaryName }}
app.kubernetes.io/managed-by: helm
app.kubernetes.io/part-of: revenu-platform
{{- end -}}

{{- define "exchangeos.image" -}}
{{- $bin := index .Values.binaries .binaryName -}}
{{- printf "%s/%s:%s" .Values.global.imageRegistry $bin.image.repository $bin.image.tag -}}
{{- end -}}
