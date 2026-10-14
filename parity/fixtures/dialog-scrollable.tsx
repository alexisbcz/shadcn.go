import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/registry/new-york-v4/ui/dialog"

// Port fixture: a scrollable page (so the scroll lock hides a real
// scrollbar and compensates its width) and a dialog whose body scrolls
// (wheel inside the content scrolls it, wheel outside is blocked).
export default function DialogScrollable() {
  return (
    <div className="flex h-[1400px] flex-col items-center gap-4 pt-10">
      <Dialog>
        <DialogTrigger asChild>
          <Button variant="outline">Terms</Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Terms of service</DialogTitle>
            <DialogDescription>Scroll to read everything.</DialogDescription>
          </DialogHeader>
          <div className="-mx-6 max-h-[240px] overflow-y-auto px-6 text-sm">
            {Array.from({ length: 12 }).map((_, i) => (
              <p key={i} className="mb-4 leading-normal">
                {`Section ${i + 1}. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`}
              </p>
            ))}
          </div>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">Decline</Button>
            </DialogClose>
            <Button>Accept</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      <p className="text-sm text-muted-foreground">The page scrolls.</p>
    </div>
  )
}
