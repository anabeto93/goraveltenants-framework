package exceptions

import "fmt"

type ModelNotSyncMasterException struct {
    class string
}

func NewModelNotSyncMasterException(class string) *ModelNotSyncMasterException {
    return &ModelNotSyncMasterException{class: class}
}

func (e *ModelNotSyncMasterException) Error() string {
    return fmt.Sprintf("Model of %s class is not a SyncMaster model. Make sure you're using the central model to make changes to synced resources when you're in the central context", e.class)
}
