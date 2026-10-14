import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/registry/new-york-v4/ui/accordion"

// Port fixture: a single, non-collapsible accordion with a disabled item,
// a horizontal one in a right-to-left context, and a fully disabled one.
export default function AccordionStates() {
  return (
    <div className="flex w-full max-w-md flex-col gap-6">
      <Accordion type="single" defaultValue="one">
        <AccordionItem value="one">
          <AccordionTrigger>One</AccordionTrigger>
          <AccordionContent>First content.</AccordionContent>
        </AccordionItem>
        <AccordionItem value="two" disabled>
          <AccordionTrigger>Two (disabled)</AccordionTrigger>
          <AccordionContent>Second content.</AccordionContent>
        </AccordionItem>
        <AccordionItem value="three">
          <AccordionTrigger>Three</AccordionTrigger>
          <AccordionContent>Third content.</AccordionContent>
        </AccordionItem>
      </Accordion>
      <div dir="rtl">
        <Accordion type="single" collapsible orientation="horizontal" dir="rtl" className="flex gap-4">
          <AccordionItem value="x">
            <AccordionTrigger>X</AccordionTrigger>
            <AccordionContent>X content.</AccordionContent>
          </AccordionItem>
          <AccordionItem value="y">
            <AccordionTrigger>Y</AccordionTrigger>
            <AccordionContent>Y content.</AccordionContent>
          </AccordionItem>
          <AccordionItem value="z">
            <AccordionTrigger>Z</AccordionTrigger>
            <AccordionContent>Z content.</AccordionContent>
          </AccordionItem>
        </Accordion>
      </div>
      <Accordion type="single" disabled defaultValue="p">
        <AccordionItem value="p">
          <AccordionTrigger>Disabled open</AccordionTrigger>
          <AccordionContent>Disabled accordion content.</AccordionContent>
        </AccordionItem>
        <AccordionItem value="q">
          <AccordionTrigger>Disabled closed</AccordionTrigger>
          <AccordionContent>Hidden.</AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>
  )
}
