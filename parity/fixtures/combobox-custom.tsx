"use client"

import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/registry/new-york-v4/ui/combobox"
import {
  Item,
  ItemContent,
  ItemDescription,
  ItemTitle,
} from "@/registry/new-york-v4/ui/item"

const countries = [
  { code: "ar", value: "argentina", label: "Argentina", continent: "South America" },
  { code: "au", value: "australia", label: "Australia", continent: "Oceania" },
  { code: "br", value: "brazil", label: "Brazil", continent: "South America" },
  { code: "ca", value: "canada", label: "Canada", continent: "North America" },
  { code: "cn", value: "china", label: "China", continent: "Asia" },
  { code: "eg", value: "egypt", label: "Egypt", continent: "Africa" },
  { code: "fr", value: "france", label: "France", continent: "Europe" },
]

// Port fixture: object items converted with itemToStringValue, rendered as
// Items with a title and a description, and a named combobox in a form.
export default function ComboboxCustom() {
  return (
    <form>
      <Combobox
        name="country"
        items={countries}
        itemToStringValue={(country: (typeof countries)[number]) => country.label}
      >
        <ComboboxInput placeholder="Search countries..." />
        <ComboboxContent>
          <ComboboxEmpty>No countries found.</ComboboxEmpty>
          <ComboboxList>
            {(country) => (
              <ComboboxItem key={country.code} value={country}>
                <Item size="sm" className="p-0">
                  <ItemContent>
                    <ItemTitle className="whitespace-nowrap">{country.label}</ItemTitle>
                    <ItemDescription>
                      {country.continent} ({country.code})
                    </ItemDescription>
                  </ItemContent>
                </Item>
              </ComboboxItem>
            )}
          </ComboboxList>
        </ComboboxContent>
      </Combobox>
    </form>
  )
}
