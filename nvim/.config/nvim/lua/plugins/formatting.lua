return {
	{
		"stevearc/conform.nvim",
		opts = {
			formatters = {
				golines = {
					prepend_args = { "-m", "110" },
				},
				ruff_format = {
					prepend_args = { "--line-length", "110" },
				},
				xmlstarlet = { "fo" },
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
				xml = { "xmlstarlet" },
			},
		},
	},
}
