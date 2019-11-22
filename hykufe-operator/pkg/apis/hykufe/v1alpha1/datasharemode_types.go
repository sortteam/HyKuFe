package v1alpha1


type DataShareSpec struct {
	// NFS 모드, 데이터가 NFS에 복사되고 각 파드가 NFS를 PV로 마운트하여 사용하게 된다.
	NFSMode *NFSModeSpec `json:"nfsMode,omitempty"`
}

type NFSModeSpec struct {
	IPAddress string `json:"ipAddress"`

	Path string `json:"path,omitempty"`
}
