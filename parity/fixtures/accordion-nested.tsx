import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/registry/new-york-v4/ui/accordion"

// Port fixture: a horizontal accordion read left to right, a vertical one
// with dir="rtl" (left/right arrows ignored), and an accordion nested in
// another one's content (keys, clicks and state stay per accordion; the
// nested one remounts with its defaults when its parent reopens).
export default function AccordionNested() {
  return (
    <div className="flex w-full max-w-md flex-col gap-6">
      <Accordion type="single" collapsible orientation="horizontal" className="flex gap-4">
        <AccordionItem value="l">
          <AccordionTrigger>L</AccordionTrigger>
          <AccordionContent>L content.</AccordionContent>
        </AccordionItem>
        <AccordionItem value="m">
          <AccordionTrigger>M</AccordionTrigger>
          <AccordionContent>M content.</AccordionContent>
        </AccordionItem>
        <AccordionItem value="n">
          <AccordionTrigger>N</AccordionTrigger>
          <AccordionContent>N content.</AccordionContent>
        </AccordionItem>
      </Accordion>
      <Accordion type="multiple" dir="rtl">
        <AccordionItem value="outer-1">
          <AccordionTrigger>Outer one</AccordionTrigger>
          <AccordionContent>
            <Accordion type="single" collapsible defaultValue="inner-1">
              <AccordionItem value="inner-1">
                <AccordionTrigger>Inner one</AccordionTrigger>
                <AccordionContent>Inner one content.</AccordionContent>
              </AccordionItem>
              <AccordionItem value="inner-2">
                <AccordionTrigger>Inner two</AccordionTrigger>
                <AccordionContent>Inner two content.</AccordionContent>
              </AccordionItem>
            </Accordion>
          </AccordionContent>
        </AccordionItem>
        <AccordionItem value="outer-2">
          <AccordionTrigger>Outer two</AccordionTrigger>
          <AccordionContent>Outer two content.</AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>
  )
}
