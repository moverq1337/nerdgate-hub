import { StrictMode } from "react"
import { createRoot } from "react-dom/client"
import { BrowserRouter, Routes, Route } from "react-router-dom"
import "./index.css"
import App from "./App.tsx"
import { GetStartedEN } from "@/pages/GetStartedEN.tsx"
import { OperationsEN } from "@/pages/OperationsEN.tsx"
import { LandingRU } from "@/pages/LandingRU.tsx"
import { GetStartedRU } from "@/pages/GetStartedRU.tsx"
import { OperationsRU } from "@/pages/OperationsRU.tsx"
import { Loader } from "@/components/Loader.tsx"

const basename = import.meta.env.BASE_URL.replace(/\/$/, "")

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter basename={basename}>
      <Loader />
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/en/get-started" element={<GetStartedEN />} />
        <Route path="/en/get-started.html" element={<GetStartedEN />} />
        <Route path="/en/operations" element={<OperationsEN />} />
        <Route path="/en/operations.html" element={<OperationsEN />} />
        <Route path="/ru" element={<LandingRU />} />
        <Route path="/ru/" element={<LandingRU />} />
        <Route path="/ru/get-started" element={<GetStartedRU />} />
        <Route path="/ru/get-started.html" element={<GetStartedRU />} />
        <Route path="/ru/operations" element={<OperationsRU />} />
        <Route path="/ru/operations.html" element={<OperationsRU />} />
        <Route path="*" element={<App />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
)
