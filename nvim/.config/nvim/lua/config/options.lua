-- Options are automatically loaded before lazy.nvim startup
-- Default options that are always set: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/options.lua
-- Add any additional options here
vim.opt.shiftwidth = 4 -- the number of spaces inserted for each indentation
vim.opt.tabstop = 4 -- the number of spaces inserted for a tab

-- Add templ file extension to support go templ.
vim.filetype.add({
	extension = {
		templ = "templ",
	},
})
