import axios from "axios"

interface RegisterRequestData {
  name: string
  password: string
}

export const registerApi = async (data: RegisterRequestData) => {
  const response = await axios.post("/api/register", data)
  return response.data
}
