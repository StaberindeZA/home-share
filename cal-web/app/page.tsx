import { listEntries } from "./actions";
import UserButton from "./components/UserButton";
import { users } from "./constants";

export default async function Home() {
  const rows = await listEntries()
  return (
    <>
      <h1>Do you have a meeting?</h1>
      <table>
        <thead>
          <tr key="header">
            <th>Time</th>
            <th>{users[0].name}</th>
            <th>{users[1].name}</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={row.start}>
              <th scope="row">{row.start.slice(0, 5)}</th>
              <td><UserButton startTime={row.start} endTime={row.end} user={users[0]} entryId={row.entryIds.at(0)} /></td>
              <td><UserButton startTime={row.start} endTime={row.end} user={users[1]} entryId={row.entryIds.at(1)} /></td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
}
