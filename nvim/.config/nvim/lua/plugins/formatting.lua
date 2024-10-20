return {
	{
		"stevearc/conform.nvim",
		opts = {
			formatters = {
				golines = {
					prepend_args = { "-m", "110" },
				},
				black = {
					prepend_args = { "--line-length", "110" },
				},
			},
			formatters_by_ft = {
				python = { "isort", "black" },
				lua = { "stylua" },
				go = { "gofumpt", "goimports", "golines" },
				sh = { "shfmt" },
				proto = { "buf" },
				c = { "clang-format" },
				cpp = { "clang-format" },
				cc = { "clang-format" },
			},
		},
	},
}
