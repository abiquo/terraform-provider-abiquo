resource "abiquo_role" "test" {
  name = "test role"
  privileges = [
    "APPLIB_UPLOAD_IMAGE",
    "VAPP_CREATE_STATEFUL",
    "VDC_MANAGE_VAPP",
    "VM_ACTION_PLAN_MANAGE",
  ]
}
