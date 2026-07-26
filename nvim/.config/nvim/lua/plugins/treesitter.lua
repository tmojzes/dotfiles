return {
	-- add more treesitter parsers
	{
		"nvim-treesitter/nvim-treesitter",
		opts = {
			ensure_installed = {
				"bash",
				"html",
				"scss",
				"templ",
				"javascript",
				"json",
				"lua",
				"markdown",
				"markdown_inline",
				"python",
				"regex",
				"vim",
				"yaml",
				"go",
				"gomod",
				"gowork",
				"gosum",
				"http",
				"rust",
				"gitignore",
				"css",
				"sql",
				"xml",
				"cue",
				"oxlint",
				"oxfmt",
			},
		},
	},
	-- {
	-- 	"smrtrfszm/dataprime.nvim",
	-- 	dependencies = { "nvim-treesitter/nvim-treesitter" },
	-- 	-- Optionally:
	-- 	-- dependencies = {'nvim-treesitter/nvim-treesitter', 'numToStr/Comment.nvim'},
	-- 	config = function(_, _)
	-- 		require("dataprime").setup()
	-- 	end,
	-- },
}
