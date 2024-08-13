import { SingleStoreLogo } from "@/components/logo/singlestore";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { BACKEND_URL } from "@/consts/config";
import axios from "axios";
import { useEffect, useState } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";
import { setCity, useCity } from "@/lib/store";
import { useNavigate } from "react-router-dom";

interface HeaderProps {
  currentPage: string;
}

export default function Header({ currentPage }: HeaderProps) {
  const [cities, setCities] = useState(["San Francisco"]);
  const selectedCity = useCity();

  const navigate = useNavigate();

  const getCities = async () => {
    const response = await axios.get(`${BACKEND_URL}/cities`);
    setCities(response.data);
  };

  useEffect(() => {
    getCities();
  }, []);

  return (
    <Card className="w-full p-2">
      <div className="flex items-center justify-between gap-2">
        <SingleStoreLogo size={40} />
        <div className="flex items-center gap-2">
          <Button
            variant={currentPage === "dashboard" ? "default" : "ghost"}
            onClick={() => navigate("/dashboard")}
          >
            Dashboard
          </Button>
          <Button
            variant={currentPage === "analytics" ? "default" : "ghost"}
            onClick={() => navigate("/analytics")}
          >
            Analytics
          </Button>
        </div>
        <Select onValueChange={(value) => setCity(value)} value={selectedCity}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="City" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="All">All Cities</SelectItem>
            {cities.map((city) => (
              <SelectItem value={city}>{city}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    </Card>
  );
}
