local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
  -- bootstrap lazy.nvim
  -- stylua: ignore
  vim.fn.system({ "git", "clone", "--filter=blob:none", "https://github.com/folke/lazy.nvim.git", "--branch=stable", lazypath })
end
vim.opt.rtp:prepend(vim.env.LAZY or lazypath)

require("lazy").setup({
	spec = {
		-- add LazyVim and import its plugins
		{ "LazyVim/LazyVim", import = "lazyvim.plugins" },
		-- LazyVim language extras
		{ import = "lazyvim.plugins.extras.lang.json" },
		{ import = "lazyvim.plugins.extras.lang.go" },
		{ import = "lazyvim.plugins.extras.lang.tailwind" },
		{ import = "lazyvim.plugins.extras.lang.rust" },
		{ import = "lazyvim.plugins.extras.lang.docker" },
		{ import = "lazyvim.plugins.extras.lang.python" },
		{ import = "lazyvim.plugins.extras.lang.terraform" },
		{ import = "lazyvim.plugins.extras.lang.yaml" },
		{ import = "lazyvim.plugins.extras.lang.markdown" },
		-- LazyVim yanky extra for better clipboard support
		{ import = "lazyvim.plugins.extras.coding.yanky" },
		-- LazyVim Debug Adapter protocol extra for better debugging support
		{ import = "lazyvim.plugins.extras.dap.core" },
		-- LazyVim test extra for better testing support
		{ import = "lazyvim.plugins.extras.test.core" },
		-- LazyVim mini extras
		{ import = "lazyvim.plugins.extras.ui.mini-indentscope" },
		{ import = "lazyvim.plugins.extras.coding.mini-surround" },
		{ import = "lazyvim.plugins.extras.coding.mini-comment" },
		{ import = "lazyvim.plugins.extras.editor.mini-move" },
		-- LazyVim ui extras
		{ import = "lazyvim.plugins.extras.ui.treesitter-context" },
		-- import/override with your plugins
		{ import = "plugins" },
	},
	defaults = {
		-- By default, only LazyVim plugins will be lazy-loaded. Your custom plugins will load during startup.
		-- If you know what you're doing, you can set this to `true` to have all your custom plugins lazy-loaded by default.
		lazy = false,
		-- It's recommended to leave version=false for now, since a lot the plugin that support versioning,
		-- have outdated releases, which may break your Neovim install.
		version = false, -- always use the latest git commit
		-- version = "*", -- try installing the latest stable version for plugins that support semver
	},
	install = { colorscheme = { "tokyonight", "night" } },
	checker = { enabled = true }, -- automatically check for plugin updates
	performance = {
		rtp = {
			-- disable some rtp plugins
			disabled_plugins = {
				"gzip",
				"tarPlugin",
				"tohtml",
				"tutor",
				"zipPlugin",
			},
		},
	},
})

-- Set colorscheme
vim.cmd("colorscheme adwaita")
