import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/registry/new-york-v4/ui/dialog"

// Port fixture: a dialog opened from another dialog, to exercise layer
// stacking (Escape and outside clicks close the top one only, the lower
// content stops receiving pointer events, focus returns to the inner
// trigger).
export default function DialogNested() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">Open outer</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Outer dialog</DialogTitle>
          <DialogDescription>Open another dialog on top of this one.</DialogDescription>
        </DialogHeader>
        <Dialog>
          <DialogTrigger asChild>
            <Button>Open inner</Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-sm">
            <DialogHeader>
              <DialogTitle>Inner dialog</DialogTitle>
              <DialogDescription>Escape closes this one first.</DialogDescription>
            </DialogHeader>
            <DialogFooter showCloseButton />
          </DialogContent>
        </Dialog>
      </DialogContent>
    </Dialog>
  )
}
