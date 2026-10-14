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

// Port fixture: a non-modal dialog (no overlay, no scroll lock, the page
// stays interactive; outside pointer and focus dismiss it), without the
// corner close button and with the footer's close button.
export default function DialogNonModal() {
  return (
    <div className="flex items-center gap-2">
      <Dialog modal={false}>
        <DialogTrigger asChild>
          <Button variant="outline">Non-modal</Button>
        </DialogTrigger>
        <DialogContent showCloseButton={false}>
          <DialogHeader>
            <DialogTitle>Non-modal dialog</DialogTitle>
            <DialogDescription>
              The page behind stays interactive.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter showCloseButton>
            <Button>Save</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      <Button variant="secondary">Outside</Button>
    </div>
  )
}
