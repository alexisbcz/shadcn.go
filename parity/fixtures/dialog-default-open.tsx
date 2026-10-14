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

// Port fixture: a dialog open on load (defaultOpen). Non-modal, so the
// harness' automatic interactions can still reach the trigger.
export default function DialogDefaultOpen() {
  return (
    <Dialog defaultOpen modal={false}>
      <DialogTrigger asChild>
        <Button variant="outline">Open again</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle>Welcome</DialogTitle>
          <DialogDescription>This dialog is open when the page loads.</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button>Got it</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
