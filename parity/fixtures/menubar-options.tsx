import {
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarPortal,
  MenubarSeparator,
  MenubarSub,
  MenubarSubContent,
  MenubarSubTrigger,
  MenubarTrigger,
} from "@/registry/new-york-v4/ui/menubar"

// Port fixture: a menubar without loop and with a menu open on load
// (defaultValue), a disabled trigger, a portalled nested submenu; a
// right-to-left menubar.
export default function MenubarOptions() {
  return (
    <div className="flex flex-col items-center gap-24">
      <Menubar loop={false} defaultValue="edit">
        <MenubarMenu value="file">
          <MenubarTrigger>File</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>New</MenubarItem>
            <MenubarSub>
              <MenubarSubTrigger>Export</MenubarSubTrigger>
              <MenubarPortal>
                <MenubarSubContent>
                  <MenubarItem>PDF</MenubarItem>
                  <MenubarSub>
                    <MenubarSubTrigger>Image</MenubarSubTrigger>
                    <MenubarSubContent>
                      <MenubarItem>PNG</MenubarItem>
                      <MenubarItem>SVG</MenubarItem>
                    </MenubarSubContent>
                  </MenubarSub>
                </MenubarSubContent>
              </MenubarPortal>
            </MenubarSub>
          </MenubarContent>
        </MenubarMenu>
        <MenubarMenu value="edit">
          <MenubarTrigger>Edit</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>Undo</MenubarItem>
            <MenubarSeparator />
            <MenubarItem disabled>Redo</MenubarItem>
          </MenubarContent>
        </MenubarMenu>
        <MenubarMenu value="locked">
          <MenubarTrigger disabled>Locked</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>Never shown</MenubarItem>
          </MenubarContent>
        </MenubarMenu>
        <MenubarMenu value="help">
          <MenubarTrigger>Help</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>About</MenubarItem>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>
      <Menubar dir="rtl">
        <MenubarMenu>
          <MenubarTrigger>Right</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>One</MenubarItem>
            <MenubarSub>
              <MenubarSubTrigger>Sub</MenubarSubTrigger>
              <MenubarSubContent>
                <MenubarItem>Inner</MenubarItem>
              </MenubarSubContent>
            </MenubarSub>
          </MenubarContent>
        </MenubarMenu>
        <MenubarMenu>
          <MenubarTrigger>Left</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>Two</MenubarItem>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>
    </div>
  )
}
