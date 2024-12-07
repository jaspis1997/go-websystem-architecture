import { createSignal } from "solid-js";
import "./App.scss";
import solidLogo from "./assets/solid.svg";

function App() {
  const [count, setCount] = createSignal(0);

  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank" rel="noreferrer">
          <img src={"/static/vite.svg"} class="logo" alt="Vite logo" />
        </a>
        <a href="https://solidjs.com" target="_blank" rel="noreferrer">
          <img src={solidLogo} class="logo solid" alt="Solid logo" />
        </a>
      </div>
      <h1>Vite + Solid</h1>
      <div class="card">
        {/* biome-ignore lint/a11y/useButtonType: <explanation> */}
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count()}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p class="read-the-docs">
        Click on the Vite and Solid logos to learn more
      </p>
    </>
  );
}

export default App;
