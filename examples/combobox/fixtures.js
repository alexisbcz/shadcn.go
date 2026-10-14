// combobox-groups renders its groups with a function of the filtered
// groups: {index < timezones.length - 1 && <ComboboxSeparator />}. The Go
// fixture renders every group with its separator (or none, for the
// second combobox); after each render of the list this keeps a group's
// separator only while its index among the rendered groups is below the
// number of groups minus one. The list carries
// data-example-combobox-groups (the number of groups) until then.
document.addEventListener("shadcn:combobox-render", (e) => {
  const list = e.target
  if (!(list instanceof Element)) return
  if (list.hasAttribute("data-example-combobox-groups")) {
    list.__groupCount = Number(list.getAttribute("data-example-combobox-groups"))
    list.removeAttribute("data-example-combobox-groups")
  }
  if (list.__groupCount === undefined) return
  e.detail.groups.forEach((group, index) => {
    if (!("__separator" in group)) group.__separator = group.querySelector(':scope > [data-slot="combobox-separator"]')
    const sep = group.__separator
    if (!sep) return
    if (index < list.__groupCount - 1) {
      if (!sep.isConnected || sep.parentElement !== group) group.append(sep)
    } else sep.remove()
  })
})
