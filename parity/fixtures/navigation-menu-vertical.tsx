import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
} from "@/registry/new-york-v4/ui/navigation-menu"

// Port fixture: vertical orientation (ArrowUp/ArrowDown between triggers,
// ArrowRight enters the open panel), without the viewport.
export default function NavigationMenuVertical() {
  return (
    <NavigationMenu
      orientation="vertical"
      viewport={false}
      delayDuration={0}
    >
      <NavigationMenuList className="flex-col items-start">
        <NavigationMenuItem value="one">
          <NavigationMenuTrigger>One</NavigationMenuTrigger>
          <NavigationMenuContent className="group-data-[viewport=false]/navigation-menu:top-0 group-data-[viewport=false]/navigation-menu:mt-0 md:left-full md:ml-2">
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
          <NavigationMenuContent className="group-data-[viewport=false]/navigation-menu:top-0 group-data-[viewport=false]/navigation-menu:mt-0 md:left-full md:ml-2">
            <ul className="grid w-[240px] gap-2">
              <li>
                <NavigationMenuLink href="#c">Gamma</NavigationMenuLink>
              </li>
              <li>
                <NavigationMenuLink href="#d">Delta</NavigationMenuLink>
              </li>
              <li>
                <NavigationMenuLink href="#e">Epsilon</NavigationMenuLink>
              </li>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
        <NavigationMenuItem value="three">
          <NavigationMenuLink href="#three">Three</NavigationMenuLink>
        </NavigationMenuItem>
      </NavigationMenuList>
    </NavigationMenu>
  )
}
