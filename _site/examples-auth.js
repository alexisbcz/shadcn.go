// User auth form (examples/authentication/components/user-auth-form.tsx):
// submitting disables the form for 3s, with spinners in the buttons.
const spinner = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-loader-circle size-4 animate-spin" role="status" aria-label="Loading"><path d="M21 12a9 9 0 1 1-6.219-8.56"></path></svg>'
for (const root of document.querySelectorAll("[data-user-auth-form]")) {
  const form = root.querySelector("form")
  const submit = form.querySelector("button")
  const github = root.querySelector(':scope > button')
  const icon = github.querySelector("svg")
  form.addEventListener("submit", (e) => {
    e.preventDefault()
    if (submit.disabled) return
    for (const el of [form.querySelector("input"), submit, github]) el.disabled = true
    submit.insertAdjacentHTML("afterbegin", spinner)
    const tmp = document.createElement("template")
    tmp.innerHTML = spinner
    icon.replaceWith(tmp.content.firstChild)
    setTimeout(() => {
      for (const el of [form.querySelector("input"), submit, github]) el.disabled = false
      submit.querySelector("svg")?.remove()
      github.querySelector("svg").replaceWith(icon)
    }, 3000)
  })
}
