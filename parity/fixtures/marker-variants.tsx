import {
  BookOpenCheck,
  FileTextIcon,
  GitBranchIcon,
  RotateCcwIcon,
  SearchIcon,
} from "lucide-react"

import {
  Marker,
  MarkerContent,
  MarkerIcon,
} from "@/registry/new-york-v4/ui/marker"
import { Spinner } from "@/registry/new-york-v4/ui/spinner"

// Port fixture: every Marker variant (default implied and explicit), with
// and without icons, status markers with a spinner and the shimmer text,
// asChild as a link and a button, links inside the content, and className
// overrides.
export default function MarkerVariants() {
  return (
    <div className="grid w-full max-w-3xl grid-cols-2 gap-x-10 gap-y-6 py-8">
      <div className="flex flex-col gap-6">
        <Marker>
          <MarkerContent>A default marker for inline notes.</MarkerContent>
        </Marker>
        <Marker variant="default">
          <MarkerIcon>
            <GitBranchIcon />
          </MarkerIcon>
          <MarkerContent>Switched to a new branch</MarkerContent>
        </Marker>
        <Marker variant="separator">
          <MarkerContent>Conversation compacted</MarkerContent>
        </Marker>
        <Marker variant="separator">
          <MarkerIcon>
            <SearchIcon />
          </MarkerIcon>
          <MarkerContent>Explored 4 files</MarkerContent>
        </Marker>
        <Marker variant="border">
          <MarkerContent>A border marker for row boundaries.</MarkerContent>
        </Marker>
        <Marker variant="border">
          <MarkerIcon>
            <FileTextIcon />
          </MarkerIcon>
          <MarkerContent>Opened implementation notes</MarkerContent>
        </Marker>
        <Marker className="flex-col">
          <MarkerIcon>
            <BookOpenCheck />
          </MarkerIcon>
          <MarkerContent>Syncing completed</MarkerContent>
        </Marker>
      </div>
      <div className="flex flex-col gap-6">
        <Marker role="status">
          <MarkerIcon>
            <Spinner />
          </MarkerIcon>
          <MarkerContent>Compacting conversation</MarkerContent>
        </Marker>
        <Marker variant="separator" role="status">
          <MarkerIcon>
            <Spinner />
          </MarkerIcon>
          <MarkerContent>Running tests</MarkerContent>
        </Marker>
        <Marker asChild>
          <a href="#links-and-buttons">
            <MarkerIcon>
              <GitBranchIcon />
            </MarkerIcon>
            <MarkerContent>View the pull request</MarkerContent>
          </a>
        </Marker>
        <Marker asChild variant="border">
          <button
            type="button"
            className="transition-colors hover:text-foreground"
          >
            <MarkerIcon className="text-foreground">
              <RotateCcwIcon />
            </MarkerIcon>
            <MarkerContent>Revert this change</MarkerContent>
          </button>
        </Marker>
        <Marker>
          <MarkerContent>
            Read the <a href="#docs">documentation</a> for details.
          </MarkerContent>
        </Marker>
        <Marker variant="separator">
          <MarkerContent className="font-medium text-foreground">
            Today
          </MarkerContent>
        </Marker>
        <Marker variant="separator" className="text-xs">
          <MarkerContent>Worked for 42s</MarkerContent>
        </Marker>
      </div>
    </div>
  )
}
