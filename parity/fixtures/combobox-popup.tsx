"use client"

import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
  ComboboxTrigger,
  ComboboxValue,
} from "@/registry/new-york-v4/ui/combobox"

const countries = [
  { code: "", value: "", continent: "", label: "Select country" },
  { code: "ar", value: "argentina", label: "Argentina", continent: "South America" },
  { code: "au", value: "australia", label: "Australia", continent: "Oceania" },
  { code: "br", value: "brazil", label: "Brazil", continent: "South America" },
  { code: "ca", value: "canada", label: "Canada", continent: "North America" },
  { code: "cn", value: "china", label: "China", continent: "Asia" },
  { code: "fr", value: "france", label: "France", continent: "Europe" },
  { code: "de", value: "germany", label: "Germany", continent: "Europe" },
  { code: "jp", value: "japan", label: "Japan", continent: "Asia" },
]

// Port fixture: a trigger button showing the value (ComboboxTrigger with a
// render prop, ComboboxValue) and the input inside the popup
// (showTrigger={false}), with object items and a default value.
export default function ComboboxPopup() {
  return (
    <Combobox items={countries} defaultValue={countries[0]}>
      <ComboboxTrigger
        render={
          <Button variant="outline" className="w-64 justify-between font-normal" />
        }
      >
        <ComboboxValue />
      </ComboboxTrigger>
      <ComboboxContent>
        <ComboboxInput showTrigger={false} placeholder="Search" />
        <ComboboxEmpty>No items found.</ComboboxEmpty>
        <ComboboxList>
          {(item) => (
            <ComboboxItem key={item.code} value={item}>
              {item.label}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  )
}
