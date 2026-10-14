"use client"

import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/registry/new-york-v4/ui/combobox"

const frameworks = ["Next.js", "SvelteKit", "Nuxt.js", "Remix", "Astro"] as const

function Items() {
  return (
    <>
      <ComboboxEmpty>No items found.</ComboboxEmpty>
      <ComboboxList>
        {(item) => (
          <ComboboxItem key={item} value={item}>
            {item}
          </ComboboxItem>
        )}
      </ComboboxList>
    </>
  )
}

// Port fixture: the Combobox props the basic fixture does not use:
// autoHighlight, a default value with the clear button, a disabled input
// and an invalid one.
export default function ComboboxStates() {
  return (
    <div className="flex w-72 flex-col gap-4">
      <Combobox items={frameworks} autoHighlight>
        <ComboboxInput placeholder="Auto highlight" />
        <ComboboxContent>
          <Items />
        </ComboboxContent>
      </Combobox>
      <Combobox items={frameworks} defaultValue={frameworks[0]}>
        <ComboboxInput placeholder="With clear" showClear />
        <ComboboxContent>
          <Items />
        </ComboboxContent>
      </Combobox>
      <Combobox items={frameworks}>
        <ComboboxInput placeholder="Disabled" disabled />
        <ComboboxContent>
          <Items />
        </ComboboxContent>
      </Combobox>
      <Combobox items={frameworks}>
        <ComboboxInput placeholder="Invalid" aria-invalid="true" />
        <ComboboxContent>
          <Items />
        </ComboboxContent>
      </Combobox>
    </div>
  )
}
