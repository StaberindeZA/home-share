import { User } from "../types";

interface UserCheckboxProps {
  userOne: User,
  userTwo: User,
  startChecked: boolean,
}

export function UserCheckbox({ userOne, userTwo, startChecked }: UserCheckboxProps) {
  return (
    <div className="inline-flex items-center gap-4">
      <label htmlFor="switch-component-on" className="text-white text-lg cursor-pointer">{userOne.name}</label>

      <div className="relative inline-block w-11 h-5">
        <input id="switch-component-on" name="user-checkbox" type="checkbox" defaultChecked={startChecked} className="peer appearance-none w-11 h-5 bg-green-600 rounded-full checked:bg-blue-600 cursor-pointer transition-colors duration-300" />
        <label htmlFor="switch-component-on" className="absolute top-0 left-0 w-5 h-5 bg-white rounded-full border border-slate-300 shadow-sm transition-transform duration-300 peer-checked:translate-x-6 peer-checked:border-slate-800 cursor-pointer">
        </label>
      </div>

      <label htmlFor="switch-component-on" className="text-white text-lg cursor-pointer">{userTwo.name}</label>
    </div>
  )
}
