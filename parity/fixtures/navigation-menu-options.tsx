import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
} from "@/registry/new-york-v4/ui/navigation-menu"

// Port fixture: an item open on load (defaultValue, explicit item values),
// no hover delay, right to left, without the viewport.
export default function NavigationMenuOptions() {
  return (
    <NavigationMenu
      viewport={false}
      defaultValue="two"
      delayDuration={0}
      skipDelayDuration={0}
      dir="rtl"
    >
      <NavigationMenuList>
        <NavigationMenuItem value="one">
          <NavigationMenuTrigger>One</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="grid w-[200px] gap-2">
              <li>
                <NavigationMenuLink href="#a">Alpha</NavigationMenuLink>
              </li>
              <li>
                <NavigationMenuLink href="#b">Beta</NavigationMenuLink>
              </li>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
        <NavigationMenuItem value="two">
          <NavigationMenuTrigger>Two</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="grid w-[240px] gap-2">
              <li>
                <NavigationMenuLink href="#c">Gamma</NavigationMenuLink>
              </li>
              <li>
                <NavigationMenuLink href="#d">Delta</NavigationMenuLink>
              </li>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
        <NavigationMenuItem value="three">
          <NavigationMenuTrigger>Three</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="grid w-[200px] gap-2">
              <li>
                <NavigationMenuLink href="#e">Epsilon</NavigationMenuLink>
              </li>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
      </NavigationMenuList>
    </NavigationMenu>
  )
}
