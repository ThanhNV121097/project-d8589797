// Composition root: every story mounts its own component here by adding one
// import at the top and one element in the list below. Do not put finished
// markup here — that lives in the story's component under components/.
import DisplayHelloWordFromDatabase from "@/components/DisplayHelloWordFromDatabase";

export default function Home() {
  return (
    <main>
      <DisplayHelloWordFromDatabase />
    </main>
  );
}
