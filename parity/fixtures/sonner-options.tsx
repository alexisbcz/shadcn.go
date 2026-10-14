// The upstream app's copy of sonner (parity/node_modules has another one,
// whose toasts would never reach the upstream <Toaster />).
// The harness mounts the upstream <Toaster /> for files that import from "sonner".
import { toast } from "../../ui/apps/v4/node_modules/sonner/dist/index.mjs"

import { Button } from "@/registry/new-york-v4/ui/button"

// Port fixture: per-toast options of sonner (rich colors, close button,
// cancel button, a toast that stays, another position, dismissing, a
// rejected promise, loading, a non-dismissible toast).
export default function SonnerOptions() {
  return (
    <div className="flex flex-wrap gap-2">
      <Button
        variant="outline"
        onClick={() =>
          toast.success("Rich success", {
            richColors: true,
            closeButton: true,
            description: "With a close button",
          })
        }
      >
        Rich
      </Button>
      <Button
        variant="outline"
        onClick={() =>
          toast("Are you sure?", {
            duration: Infinity,
            cancel: { label: "Cancel", onClick: () => {} },
            action: { label: "Confirm", onClick: () => {} },
          })
        }
      >
        Confirm
      </Button>
      <Button
        variant="outline"
        onClick={() => toast.info("At the top", { position: "top-center" })}
      >
        Top
      </Button>
      <Button variant="outline" onClick={() => toast.dismiss()}>
        Dismiss all
      </Button>
      <Button
        variant="outline"
        onClick={() =>
          toast.promise(
            () =>
              new Promise((_, reject) =>
                setTimeout(() => reject(new Error("nope")), 1000)
              ),
            { loading: "Saving...", success: "Saved", error: "Could not save" }
          )
        }
      >
        Reject
      </Button>
      <Button
        variant="outline"
        onClick={() =>
          toast.warning("Cannot be swiped", {
            dismissible: false,
            invert: true,
          })
        }
      >
        Sticky
      </Button>
    </div>
  )
}
