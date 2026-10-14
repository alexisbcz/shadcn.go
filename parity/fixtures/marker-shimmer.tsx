import { Marker, MarkerContent } from "@/registry/new-york-v4/ui/marker"

// Port fixture: status markers whose text uses the shimmer animation.
export default function MarkerShimmer() {
  return (
    <div className="flex w-full max-w-sm flex-col gap-8 py-12">
      <Marker role="status">
        <MarkerContent className="shimmer">Thinking...</MarkerContent>
      </Marker>
      <Marker variant="separator" role="status">
        <MarkerContent className="shimmer">Reading 4 files</MarkerContent>
      </Marker>
    </div>
  )
}
