import { Routes, Route } from "react-router-dom";
import Layout from "./components/Layout";
import ReleasesPage from "./pages/ReleasesPage";
import ReleasePage from "./pages/ReleasePage";
import ArtistsPage from "./pages/ArtistsPage";
import ArtistPage from "./pages/ArtistPage";

export default function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<ReleasesPage />} />
        <Route path="/releases/:id" element={<ReleasePage />} />
        <Route path="/artists" element={<ArtistsPage />} />
        <Route path="/artists/:id" element={<ArtistPage />} />
      </Routes>
    </Layout>
  );
}
