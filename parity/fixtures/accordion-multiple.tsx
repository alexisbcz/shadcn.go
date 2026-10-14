import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/registry/new-york-v4/ui/accordion"

// Port fixture: type="multiple" with two items open by default and a
// disabled item.
export default function AccordionMultiple() {
  return (
    <Accordion type="multiple" className="w-full max-w-md" defaultValue={["a", "c"]}>
      <AccordionItem value="a">
        <AccordionTrigger>Is it accessible?</AccordionTrigger>
        <AccordionContent>Yes. It adheres to the WAI-ARIA design pattern.</AccordionContent>
      </AccordionItem>
      <AccordionItem value="b">
        <AccordionTrigger>Is it styled?</AccordionTrigger>
        <AccordionContent>Yes. It comes with default styles that match the other components.</AccordionContent>
      </AccordionItem>
      <AccordionItem value="c">
        <AccordionTrigger>Is it animated?</AccordionTrigger>
        <AccordionContent>Yes. It&apos;s animated by default, but you can disable it if you prefer.</AccordionContent>
      </AccordionItem>
      <AccordionItem value="d" disabled>
        <AccordionTrigger>Is it disabled?</AccordionTrigger>
        <AccordionContent>This item cannot be opened.</AccordionContent>
      </AccordionItem>
      <AccordionItem value="e">
        <AccordionTrigger>Can it nest?</AccordionTrigger>
        <AccordionContent>Yes.</AccordionContent>
      </AccordionItem>
    </Accordion>
  )
}
