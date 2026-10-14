import { HomeIcon, InboxIcon } from "lucide-react"

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarRail,
  SidebarTrigger,
} from "@/registry/new-york-v4/ui/sidebar"

// Port fixture: a right-side sidebar with the inset variant and the rail,
// starting collapsed (defaultOpen={false}).
export default function SidebarRight() {
  return (
    <SidebarProvider defaultOpen={false}>
      <SidebarInset>
        <header className="flex h-12 items-center justify-end px-4">
          <SidebarTrigger />
        </header>
        <div className="p-4 text-sm">Main content</div>
      </SidebarInset>
      <Sidebar side="right" variant="inset">
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Mail</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton>
                    <HomeIcon />
                    <span>Home</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton>
                    <InboxIcon />
                    <span>Inbox</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarRail />
      </Sidebar>
    </SidebarProvider>
  )
}
