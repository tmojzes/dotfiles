return {
	-- tools
	{
		"williamboman/mason.nvim",
		opts = function(_, opts)
			vim.list_extend(opts.ensure_installed, {
				"stylua",
				"shellcheck",
				"shfmt",
				"flake8",
				"luacheck",
				"tailwindcss-language-server",
				"css-lsp",
				"html-lsp",
				"htmx-lsp",
				"gopls",
				"gofumpt",
				"goimports",
				"zls",
				"codespell",
				"templ",
				"helm-ls",
				"bash-language-server",
				"rnix-lsp",
				"robotframework-lsp",
				"staticcheck",
				"buf-language-server",
				"buf",
				"ols",
			})
		end,
	},

	-- the opts function can also be used to change the default opts:
	{
		"nvim-lualine/lualine.nvim",
		event = "VeryLazy",
		opts = function(_, opts)
			table.insert(opts.sections.lualine_x, "😄")
		end,
	},

	-- lsp servers
	{
		"neovim/nvim-lspconfig",
		--FIXME: LSP does not work if sourcekit is not installed.
		-- config = function()
		-- 	local lspconfig = require("lspconfig")
		-- 	lspconfig.sourcekit.setup({})
		-- end,
		opts = {
			inlay_hints = { enabled = false },
			servers = {
				cssls = {
					filetypes = { "css", "scss", "less" },
				},
				htmx = {
					filetypes = { "html", "templ" },
				},
				gopls = {
					settings = {
						gopls = {
							-- Sets go build tags, so gopls loads packges behind these build tags as well.
							buildFlags = { "-tags=integration global" },
						},
					},
				},
				tailwindcss = {
					root_dir = function(...)
						return require("lspconfig.util").root_pattern(".git")(...)
					end,
					filetypes = {
						"aspnetcorerazor",
						"astro",
						"astro-markdown",
						"blade",
						"clojure",
						"django-html",
						"htmldjango",
						"edge",
						"eelixir",
						"elixir",
						"ejs",
						"erb",
						"eruby",
						"gohtml",
						"gohtmltmpl",
						"haml",
						"handlebars",
						"hbs",
						"html",
						"html-eex",
						"heex",
						"jade",
						"leaf",
						"liquid",
						"markdown",
						"mdx",
						"mustache",
						"njk",
						"nunjucks",
						"php",
						"razor",
						"slim",
						"twig",
						"css",
						"less",
						"postcss",
						"sass",
						"scss",
						"stylus",
						"sugarss",
						"javascript",
						"javascriptreact",
						"reason",
						"rescript",
						"typescript",
						"typescriptreact",
						"vue",
						"svelte",
						"templ",
					},
				},
			},
		},
	},
}
