return {
	{
		"projekt0n/github-nvim-theme",
		lazy = false, -- make sure we load this during startup if it is your main colorscheme
		priority = 1000, -- make sure to load this before all the other start plugins
		config = function()
			require("github-theme").setup({
				-- Put the settings here
			})
		end,
	},
	{
		"catppuccin/nvim",
		name = "catppuccin",
		lazy = false,
		priority = 1000,
		config = function()
			require("catppuccin").setup({
				-- Put the settings here
			})
		end,
	},
	{
		"Mofiqul/adwaita.nvim",
		lazy = false,
		priority = 1000,
		-- configure and set on startup
		config = function()
			vim.g.adwaita_darker = false -- for darker version
			vim.g.adwaita_disable_cursorline = false -- to disable cursorline
			vim.g.adwaita_transparent = false -- makes the background transparent
		end,
	},
	{
		"rebelot/kanagawa.nvim",
		lazy = false,
		priority = 1000,
		config = function()
			require("kanagawa").setup({
				-- Put the settings here
			})
		end,
	},
	{
		"Mofiqul/vscode.nvim",
		lazy = false,
		priority = 1000,
		config = function()
			require("vscode").setup({
				-- Put the settings here
			})
		end,
	},
	{
		"ellisonleao/gruvbox.nvim",
		priority = 1000,
		config = true,
		opts = ...,
	},
}
