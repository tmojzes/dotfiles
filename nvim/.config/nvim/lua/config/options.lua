-- Options are automatically loaded before lazy.nvim startup
-- Default options that are always set: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/options.lua
-- Add any additional options here
vim.opt.shiftwidth = 4 -- the number of spaces inserted for each indentation
vim.opt.tabstop = 4 -- the number of spaces inserted for a tab
vim.opt.colorcolumn = "110" -- sets a vertical line length marker
vim.opt.number = true
vim.opt.relativenumber = false
vim.g.lazyvim_python_ruff = "ruff"
vim.g.lazyvim_python_lsp = "ruff"

-- Add templ file extension to support go templ.
vim.filetype.add({
	extension = {
		templ = "templ",
	},
})

-- Add filetypes for ansible.
vim.filetype.add({
	pattern = {
		-- From: https://github.com/mfussenegger/nvim-ansible/blob/bba61168b7aef735e7f950fdfece5ef6c388eacf/ftdetect/ansible.lua#L4-L13
		[".*/defaults/.*%.ya?ml"] = "yaml.ansible",
		[".*/host_vars/.*%.ya?ml"] = "yaml.ansible",
		[".*/group_vars/.*%.ya?ml"] = "yaml.ansible",
		[".*/group_vars/.*/.*%.ya?ml"] = "yaml.ansible",
		[".*/playbook.*%.ya?ml"] = "yaml.ansible",
		[".*/playbooks/.*%.ya?ml"] = "yaml.ansible",
		[".*/roles/.*/tasks/.*%.ya?ml"] = "yaml.ansible",
		[".*/roles/.*/handlers/.*%.ya?ml"] = "yaml.ansible",
		[".*/tasks/.*%.ya?ml"] = "yaml.ansible",
		[".*/molecule/.*%.ya?ml"] = "yaml.ansible",
	},
})
