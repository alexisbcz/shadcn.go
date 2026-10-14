import {
  NativeSelect,
  NativeSelectOption,
} from "@/registry/new-york-v4/ui/native-select"

// Port fixture: the "sm" size and a className on the select.
export default function NativeSelectSm() {
  return (
    <NativeSelect size="sm" className="w-40">
      <NativeSelectOption value="">Select size</NativeSelectOption>
      <NativeSelectOption value="s">Small</NativeSelectOption>
      <NativeSelectOption value="m">Medium</NativeSelectOption>
    </NativeSelect>
  )
}
