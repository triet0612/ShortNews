import {api_url} from "./index.js"
import { get } from "svelte/store"

export class Article {
  /**
   * @param {string} id 
   * @param {string} link 
   * @param {string} title 
   * @param {Date} pubdate 
   * @param {string} pubid 
   * @param {string} summary 
   */
  constructor(id, link, title, pubdate, pubid, summary) {
    this.id = id
    this.link = link
    this.title = title
    this.pubdate = pubdate
    this.pubid = pubid
    this.summary = summary
  }
}
/**
 * @param {number} page 
 * @param {string} pub
 * @returns {Promise<Article[]>}
 */
export async function articleFromApi(page, pub) {
  /**
   * @type {Array<any>}
   */
  let jsResponse = await fetch(get(api_url) + `/articles?limit=5&start=${page*5}&summary=true&PublisherID=${pub}`)
    .then(res => res.json())
    .catch(err => {console.log(err);return []})
  /**
   * @type {Article[]}
   */
  let articles = []
  jsResponse.map(v => {
    articles.push(new Article(
      v["ArticleID"], v["Link"], v["Title"], v["PubDate"], v["PublisherID"], v["Summary"]
    ))
  })
  return articles
}
/**
 * @returns {Promise<Article>}
 */
export async function randomArticle() {
  /**
   * @type {any}
   */
  let js = await fetch(get(api_url) + `/articles/random`)
    .then(res => res.json())
    .catch(err => {console.log(err);return undefined})
  let ans = new Article(js["ArticleID"], js["Link"], js["Title"], js["PubDate"], js["PublisherID"], js["Summary"])
  return ans
}
