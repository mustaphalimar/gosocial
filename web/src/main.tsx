import { createRoot } from "react-dom/client";
import "./input.css";
import App from "./App";
import { createBrowserRouter } from "react-router-dom";
import ConfirmationPage from "./pages/confirmation-page";
import { RouterProvider } from "react-router-dom";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
  },
  {
    path: "/confirm/:token",
    element: <ConfirmationPage />,
  },
]);

createRoot(document.getElementById("root")!).render(
  <>
    <RouterProvider router={router} />
  </>,
);
