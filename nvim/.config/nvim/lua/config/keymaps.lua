-- Keymaps are automatically loaded on the VeryLazy event
-- Default keymaps that are always set: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/keymaps.lua
-- Add any additional keymaps here
local keymap = vim.keymap
local opts = { noremap = true, silent = true }

--- Set custom keymaps here.

-- In visual mode ('x'), map 'p' to 'P' to paste without overwriting the register.
vim.keymap.set("x", "p", "P", { desc = "Paste without yanking selected text" })
