package repository

const getUserById = "SELECT id, full_name, phone_number, password, successful_login FROM users WHERE id = $1"
const getUserByPhoneNumber = "SELECT id, full_name, phone_number, password, successful_login FROM users WHERE phone_number = $1"
const createUser = "INSERT INTO users (full_name, phone_number, password) VALUES ($1, $2, $3) RETURNING id"
const updateUser = "UPDATE users SET full_name = $1, phone_number = $2 WHERE id = $3"
const succesLoginIncrement = "UPDATE users SET successful_login = $1 WHERE id = $2"
