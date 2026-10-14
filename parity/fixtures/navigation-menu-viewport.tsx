import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuIndicator,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  navigationMenuTriggerStyle,
} from "@/registry/new-york-v4/ui/navigation-menu"

// Port fixture: the shared viewport (shadcn's default), the indicator, an
// active link and a disabled trigger. No hover delay, so the automatic
// hover capture is not racing the open timer.
export default function NavigationMenuViewportFixture() {
  return (
    <NavigationMenu delayDuration={0}>
      <NavigationMenuList>
        <NavigationMenuItem>
          <NavigationMenuTrigger>Getting started</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="grid w-[400px] gap-2">
              <li>
                <NavigationMenuLink href="#intro" active>
                  Introduction
                </NavigationMenuLink>
              </li>
              <li>
                <NavigationMenuLink href="#install">Installation</NavigationMenuLink>
              </li>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuTrigger>Components</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="grid w-[260px] gap-2">
              <li>
                <NavigationMenuLink href="#button">Button</NavigationMenuLink>
              </li>
              <li>
                <NavigationMenuLink href="#card">Card</NavigationMenuLink>
              </li>
              <li>
                <NavigationMenuLink href="#dialog">Dialog</NavigationMenuLink>
              </li>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuTrigger disabled>Disabled</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="w-[200px]">
              <li>
                <NavigationMenuLink href="#never">Never</NavigationMenuLink>
              </li>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuLink href="#docs" className={navigationMenuTriggerStyle()}>
            Docs
          </NavigationMenuLink>
        </NavigationMenuItem>
        <NavigationMenuIndicator />
      </NavigationMenuList>
    </NavigationMenu>
  )
}
