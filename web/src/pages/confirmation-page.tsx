import { API_URL } from "@/lib/constants";
import { Button } from "@/components/ui/button";
import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";

export default function ConfirmationPage() {
  const { token = "" } = useParams();
  const redirect = useNavigate();

  const handleConfirm = async () => {
    const response = await fetch(`${API_URL}/users/activate/${token}`, {
      method: "PUT",
    });

    if (response.ok) {
      // redirect home
      redirect("/");
    } else {
      alert("Failed to confirm");
    }
  };

  return (
    <div>
      <h1>Confirmation Page</h1>
      <Button onClick={handleConfirm}>Click to confirm</Button>
    </div>
  );
}
